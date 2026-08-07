import os
import sys

from PIL import Image, ImageFilter, ImageOps

dracula_theme = {
    "bg": (40, 42, 54),
    "mid": (68, 71, 90),
    "alt_mid": (98, 114, 164),
    "fg": (248, 248, 242),
    "cyan": (139, 233, 253),
    "green": (80, 250, 123),
    "orange": (255, 184, 108),
    "pink": (255, 121, 198),
    "purple": (189, 147, 249),
    "red": (255, 85, 85),
    "yellow": (241, 250, 140),
}


def apply_dracula(image):
    rgb_image = image.convert("RGB")
    rgb_pixels = rgb_image.load()

    grayscale_image = image.convert("L")
    bnw_image = ImageOps.colorize(
        grayscale_image,
        dracula_theme["bg"],
        dracula_theme["fg"],
        dracula_theme["mid"],
        blackpoint=10,
    )

    red_image = ImageOps.colorize(
        grayscale_image,
        dracula_theme["pink"],
        dracula_theme["red"],
        dracula_theme["red"],
        blackpoint=10,
    )

    green_image = ImageOps.colorize(
        grayscale_image,
        dracula_theme["orange"],
        dracula_theme["green"],
        dracula_theme["green"],
        blackpoint=10,
    )

    blue_image = ImageOps.colorize(
        grayscale_image,
        dracula_theme["purple"],
        dracula_theme["cyan"],
        dracula_theme["cyan"],
        blackpoint=10,
    )

    width, height = bnw_image.size

    red_mask = Image.new("L", (width, height), 0)
    red_mask_pixels = red_mask.load()
    green_mask = Image.new("L", (width, height), 0)
    green_mask_pixels = green_mask.load()
    blue_mask = Image.new("L", (width, height), 0)
    blue_mask_pixels = blue_mask.load()

    for y in range(height):
        for x in range(width):
            r, g, b = rgb_pixels[x, y]

            red_dominance = min(r - g, r - b)
            green_dominance = min(g - r, g - b)
            blue_dominance = min(b - r, b - g)

            if red_dominance >= 20 and r >= 30:
                red_mask_value = int(red_dominance * 1.2)
                red_mask_pixels[x, y] = min(red_mask_value, 255)

            if green_dominance >= 20 and g >= 30:
                green_mask_value = int(green_dominance * 1.5)
                green_mask_pixels[x, y] = min(green_mask_value, 255)

            if blue_dominance >= 20 and b >= 30:
                blue_mask_value = int(blue_dominance * 1.2)
                blue_mask_pixels[x, y] = min(blue_mask_value, 255)

    red_mask = red_mask.filter(ImageFilter.GaussianBlur(3))
    green_mask = green_mask.filter(ImageFilter.GaussianBlur(3))
    blue_mask = blue_mask.filter(ImageFilter.GaussianBlur(3))

    result = Image.composite(red_image, bnw_image, red_mask)
    result = Image.composite(green_image, result, green_mask)
    result = Image.composite(blue_image, result, blue_mask)

    return result


def apply_long(img):
    width, height = img.size

    if width < height // 9 * width:
        long_image = Image.new("RGB", (height // 9 * 16, height), dracula_theme["bg"])
    else:
        long_image = Image.new("RGB", (width, width // 16 * 9), dracula_theme["bg"])
    long_image_pixels = long_image.load()

    long_width, long_height = long_image.size

    image_pixels = img.load()

    startx, starty = (long_width - width) // 2, (long_height - height) // 2

    for y in range(height):
        for x in range(width):
            r, g, b = image_pixels[x, y]
            long_image_pixels[startx + x, starty + y] = (r, g, b)

    return long_image


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <path-to-image> [path-to-save-location]")
        sys.exit(1)
    input_path = sys.argv[1]
    try:
        output_path = sys.argv[2]
    except IndexError:
        output_path = "."

    if not os.path.exists(input_path):
        print(f"Error: {input_path} not found.")
        sys.exit(1)

    if not os.path.exists(output_path) and not output_path.endswith((".jpg", ".png")):
        os.makedirs(output_path)

    if os.path.isdir(output_path):
        base_name = os.path.basename(input_path)
        output_path = os.path.join(output_path, f"igdt_{base_name}")

    try:
        image = Image.open(input_path)
        img = apply_dracula(image)
        result = apply_dracula(image)
        # result = apply_long(img)
        result.save(output_path, "PNG")
        print(f"Image saved to {output_path} successfully.")
    except Exception as error:
        print(f"Error: {error}")
        import traceback

        traceback.print_exc()
