import os
import requests
from svglib.svglib import svg2rlg
from reportlab.graphics import renderPM

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

            # Download the SVG file
            print(f"Downloading: {url}")
            response = requests.get(url)
            response.raise_for_status()
            with open(svg_path, "wb") as svg_file:
                svg_file.write(response.content)

            # Convert the SVG to PNG
            print(f"Converting to PNG: {png_filename}")
            drawing = svg2rlg(svg_path)
            renderPM.drawToFile(drawing, png_path, fmt="PNG")
            print(f"Saved: {png_path}")
        except Exception as e:
            print(f"Failed to process {url}: {e}")

if __name__ == "__main__":
    # List of URLs to download SVGs from
    svg_urls = [
        "https://beta-cdn.wynncraft.com/nextgen/banners/CIRCLE_MIDDLE.svg",
        "https://beta-cdn.wynncraft.com/nextgen/banners/STRIPE_VERTICAL.svg",
        "https://beta-cdn.wynncraft.com/nextgen/banners/CROSS.svg"
    ]

    # Output directory for PNG files
    output_directory = "pythonhelper/output_pngs"

    # Run the download and conversion process
    download_and_convert_svgs(svg_urls, output_directory)
    print("All SVGs downloaded and converted to PNG!")
