from PIL import Image
import os

def upscale_image(input_path, output_path, scale=4):
    """
    Upscale a PNG image by the specified scale.

    :param input_path: Path to the input PNG file
    :param output_path: Path to save the upscaled image
    :param scale: Scaling factor (default is 4)
    """
    with Image.open(input_path) as img:
        # Ensure the image is in RGB or RGBA mode
        img = img.convert("RGBA") if img.mode != "RGBA" else img

        # Calculate new size
        new_size = (img.width * scale, img.height * scale)

        # Resize using nearest-neighbor scaling
        upscaled_img = img.resize(new_size, Image.NEAREST)

        # Save the upscaled image
        upscaled_img.save(output_path, "PNG")

def upscale_all_pngs(input_directory, output_directory, scale=4):
    """
    Upscale all PNG images in a directory by the specified scale.

    :param input_directory: Path to the directory containing PNG images
    :param output_directory: Path to save the upscaled images
    :param scale: Scaling factor (default is 4)
    """
    if not os.path.exists(output_directory):
        os.makedirs(output_directory)

    for file_name in os.listdir(input_directory):
        if file_name.lower().endswith(".png"):
            input_path = os.path.join(input_directory, file_name)
            output_path = os.path.join(output_directory, file_name)

            print(f"Upscaling {file_name}...")
            upscale_image(input_path, output_path, scale)

    print("Upscaling complete!")

if __name__ == "__main__":
    # Set input and output directories
    input_dir = "statscard/ranks"  # Replace with your input directory
    output_dir = "statscard/ranks_upscale"  # Replace with your output directory

    # Run the upscaling function
    upscale_all_pngs(input_dir, output_dir)
