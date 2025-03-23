from PIL import Image

# Load the uploaded sprite sheet
image_path = "/mnt/data/warrior-main-2.png"
sprite_sheet = Image.open(image_path)

# Get image dimensions
width, height = sprite_sheet.size

# Define the number of columns and rows in the sprite sheet
cols, rows = 8, 10  # Assuming 8 columns and 10 rows

# Calculate the width and height of each sprite
sprite_width = width // cols
sprite_height = height // rows

# Extract frames from 3rd and 4th row and flip them horizontally
flipped_sprites = []
for row in [2, 3]:  # 0-based index, so 3rd row is index 2, 4th row is index 3
    for col in range(cols):
        # Crop the sprite
        left = col * sprite_width
        top = row * sprite_height
        right = left + sprite_width
        bottom = top + sprite_height
        sprite = sprite_sheet.crop((left, top, right, bottom))
        
        # Flip the sprite horizontally
        flipped_sprite = sprite.transpose(Image.FLIP_LEFT_RIGHT)
        flipped_sprites.append(flipped_sprite)

# Create a new image to display the flipped sprites
output_image = Image.new("RGBA", (sprite_width * cols, sprite_height * 2))

# Paste the flipped sprites into the new image
for index, flipped_sprite in enumerate(flipped_sprites):
    x_offset = (index % cols) * sprite_width
    y_offset = (index // cols) * sprite_height
    output_image.paste(flipped_sprite, (x_offset, y_offset))

# Save the output image
flipped_image_path = "/mnt/data/flipped_sprites.png"
output_image.save(flipped_image_path)
