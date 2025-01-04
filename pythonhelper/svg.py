import os
import requests
from svglib.svglib import svg2rlg
from reportlab.graphics import renderPM
import cairosvg

# if shit doesn't work, run
"""
export PKG_CONFIG_PATH="/opt/homebrew/lib/pkgconfig:$PKG_CONFIG_PATH"
export DYLD_LIBRARY_PATH="/opt/homebrew/lib:$DYLD_LIBRARY_PATH"
export LIBRARY_PATH="/opt/homebrew/lib:$LIBRARY_PATH"
"""

def download_and_convert_svgs(svg_urls, output_dir):
    """
    Downloads SVG files from URLs and converts them to PNG files.

    Args:
        svg_urls (list): List of URLs pointing to SVG files.
        output_dir (str): Path to the directory where PNG files will be saved.
    """
    # Ensure the output directory exists
    os.makedirs(output_dir, exist_ok=True)

    for url in svg_urls:
        try:
            # Extract the filename from the URL
            svg_filename = os.path.basename(url)
            png_filename = svg_filename.replace(".svg", ".png")
            
            svg_path = os.path.join(output_dir, svg_filename)
            png_path = os.path.join(output_dir, png_filename)

            # Skip downloading and conversion if the PNG already exists
            if os.path.exists(png_path):
                print(f"Skipping: {png_filename} already exists.")
                continue

            # Download the SVG file if it doesn't already exist
            if not os.path.exists(svg_path):
                print(f"Downloading: {url}")
                response = requests.get(url)
                response.raise_for_status()
                with open(svg_path, "wb") as svg_file:
                    svg_file.write(response.content)

            # Convert the SVG to PNG with alpha support
            print(f"Converting to PNG: {png_filename}")
            cairosvg.svg2png(url=svg_path, write_to=png_path)
            print(f"Saved: {png_path}")

            # Delete the SVG file after conversion
            os.remove(svg_path)
            print(f"Deleted: {svg_path}")

            print(f"Saved: {png_path}")

        except Exception as e:
            print(f"Failed to process {url}: {e}")

if __name__ == "__main__":
    # List of URLs to download SVGs from
    # patterns = [
    #     "SQUARE_BOTTOM_LEFT", # bottom left quad
    #     "SQUARE_BOTTOM_RIGHT", # bottom right quad
    #     "SQUARE_TOP_LEFT", # top left quad
    #     "SQUARE_TOP_RIGHT", # top right quad
    #     "STRIPE_BOTTOM", # bottom stripe
    #     "STRIPE_TOP", # top stripe
    #     "STRIPE_LEFT", # left stripe
    #     "STRIPE_RIGHT", # right stripe
    #     "STRIPE_CENTER", # vertical stripe

    #     "STRIPE_MIDDLE", # horizontal stripe
    #     "STRIPE_DOWNRIGHT", # major diagonal
    #     "STRIPE_DOWNLEFT", # minor diagonal
    #     "STRIPE_SMALL", # river
    #     "CROSS", # cross
    #     "TRIANGLE_BOTTOM", # bottom triangle
    #     "TRIANGLE_TOP", # top triangle
    #     "TRIANGLES_BOTTOM", # bottom triangles
    #     "TRIANGLES_TOP", # top triangles

    #     "DIAGONAL_LEFT", # top left
    #     "DIAGONAL_RIGHT", # bottom right
    #     "DIAGONAL_LEFT_MIRROR", # bottom left
    #     "DIAGONAL_RIGHT_MIRROR", # top right
    #     "HALF_VERTICAL", # left half
    #     "HALF_HORIZONTAL", # top half
    #     "HALF_VERTICAL_MIRROR", # right half
    #     "HALF_HORIZONTAL_MIRROR", # bottom half
    #     "BORDER", # border

    #     "CREEPER", # creeper
    #     "GRADIENT", # gradient down
    #     "BRICKS", # bricks
    #     "SKULL", # skull
    #     "MOJANG", # thing

    #     "CURLY_BORDER", # triangle border
    #     "STRAIGHT_CROSS", # horizontal cross
    #     "GRADIENT_UP", # gradient up
    #     "FLOWER", # flower
    #     "RHOMBUS_MIDDLE", # rhombus
    #     "CIRCLE_MIDDLE", # circle
        
    # ]

    # svg_urls = [
    #     f"https://beta-cdn.wynncraft.com/nextgen/banners/{pattern}.svg" for pattern in patterns if pattern != ""
    # ]

    ranks = [
        "vip",
        "vipplus",
        "hero",
        "champion",
        "media",
        "item",
        "moderator",
        "administrator",
        "builder",
        "gamemaster",
        "cmd",
        "hybrid",
        "qa",
        "art",
        "music",
    ]

    svg_urls = [
        f"https://cdn.wynncraft.com/nextgen/badges/rank_{rank}.svg" for rank in ranks
    ]

    # Output directory for PNG files
    output_directory = "statscard/ranks"

    # Run the download and conversion process
    download_and_convert_svgs(svg_urls, output_directory)
    print("All SVGs downloaded and converted to PNG!")

