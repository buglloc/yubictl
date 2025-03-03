#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>

#include "pico/stdlib.h"
#include "pico/util/queue.h"
#include "pico/multicore.h"

#include "hardware/gpio.h"
#include "hardware/pio.h"
#include "hardware/clocks.h"

#include "ws2812.pio.h"


#define MIN_PIN          3
#define MAX_PIN          29
#define LED_PIN          16
#define LED_COLOR        urgb_u32(0, 0, 0x20)
#define TOUCH_TIMEOUT_MS 500
#define TOUCH_QUEUE_SIZE 10


static PIO led_pio;
static uint led_sm;
static queue_t touch_queue;

static inline uint32_t urgb_u32(uint8_t r, uint8_t g, uint8_t b) {
  return
    ((uint32_t) (r) << 8) |
    ((uint32_t) (g) << 16) |
    (uint32_t) (b);
}

void led_init() {
  uint offset;
  bool success = pio_claim_free_sm_and_add_program_for_gpio_range(&ws2812_program, &led_pio, &led_sm, &offset, LED_PIN, 1, true);
  hard_assert(success);

  ws2812_program_init(led_pio, led_sm, offset, LED_PIN, 800000, false);
}

static inline void led_on() {
  pio_sm_put_blocking(led_pio, led_sm, LED_COLOR << 8u);
}

static inline void led_off() {
  pio_sm_put_blocking(led_pio, led_sm, 0);
}

void touch_loop() {
  uint gpio;
  while (1) {
    queue_remove_blocking(&touch_queue, &gpio);

    led_on();

    gpio_put(gpio, 1);
    sleep_ms(TOUCH_TIMEOUT_MS);
    gpio_put(gpio, 0);
    sleep_ms(TOUCH_TIMEOUT_MS / 2);

    led_off();
  }
}

static inline bool is_pin_available(int gpio) {
  return (
    gpio >= MIN_PIN &&
    gpio <= MAX_PIN &&
    gpio != LED_PIN
  );
}

void printhelp() {
  puts("\nCommands:");
  puts("t<pin>\t: touch pin");
}

int main() {
  led_init();
  stdio_usb_init();

  for (uint p = MIN_PIN; p < MAX_PIN; ++p) {
    if (p == LED_PIN) {
      continue;
    }

    gpio_init(p);
    gpio_put(p, 0);
    gpio_set_dir(p, GPIO_OUT);
  }

  queue_init(&touch_queue, sizeof(uint), TOUCH_QUEUE_SIZE);
  multicore_launch_core1(touch_loop);

  while (1) {
    char c = getchar();
    switch (c) {
      case 't':
        c = getchar();
        uint pin = c;
        if (!is_pin_available(pin)) {
          printf("\ninvalid pin: %d\n", pin);
          break;
        }

        queue_add_blocking(&touch_queue, &pin);
        break;
      case '\n':
      case '\r':
        break;
      case 'h':
        printhelp();
        break;
      default:
        printf("\nUnrecognised command: %c\n", c);
        printhelp();
        break;
    }
  }
}
