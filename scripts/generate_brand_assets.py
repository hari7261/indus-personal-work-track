from pathlib import Path

from PIL import Image, ImageDraw, ImageFont


ROOT = Path(__file__).resolve().parents[1]
BUILD_DIR = ROOT / "backend" / "build"
WINDOWS_DIR = BUILD_DIR / "windows"
DOC_ASSETS = ROOT / "docs" / "assets"


def ensure_dirs() -> None:
    WINDOWS_DIR.mkdir(parents=True, exist_ok=True)
    DOC_ASSETS.mkdir(parents=True, exist_ok=True)


def generate_app_icon() -> None:
    size = 1024
    img = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # Transparent background with only foreground artwork.
    draw.rounded_rectangle((64, 64, 960, 960), radius=190, fill=(14, 165, 233, 245))
    draw.rounded_rectangle((100, 100, 924, 924), radius=170, outline=(255, 255, 255, 230), width=16)

    # Stylized board/check shape.
    draw.rounded_rectangle((220, 250, 804, 820), radius=95, fill=(255, 255, 255, 238))
    draw.rounded_rectangle((280, 348, 744, 438), radius=28, fill=(14, 165, 233, 220))
    draw.rounded_rectangle((280, 488, 558, 568), radius=24, fill=(16, 185, 129, 220))
    draw.rounded_rectangle((586, 488, 744, 568), radius=24, fill=(99, 102, 241, 215))
    draw.rounded_rectangle((280, 612, 744, 692), radius=24, fill=(249, 115, 22, 215))

    # "IP" initials for "Indus Personal".
    font = ImageFont.load_default(size=140)
    draw.text((360, 142), "IP", fill=(255, 255, 255, 250), font=font)

    appicon_path = BUILD_DIR / "appicon.png"
    img.save(appicon_path, format="PNG")

    ico_path = WINDOWS_DIR / "icon.ico"
    img.save(
        ico_path,
        format="ICO",
        sizes=[(16, 16), (24, 24), (32, 32), (48, 48), (64, 64), (128, 128), (256, 256)],
    )


def generate_repo_banner() -> None:
    width, height = 1280, 640
    img = Image.new("RGBA", (width, height), (10, 15, 27, 255))
    draw = ImageDraw.Draw(img)

    # Accent shapes.
    draw.ellipse((-140, 120, 320, 580), fill=(14, 165, 233, 180))
    draw.ellipse((920, -120, 1450, 410), fill=(99, 102, 241, 165))
    draw.rounded_rectangle((140, 110, 1140, 530), radius=40, fill=(255, 255, 255, 232))
    draw.rounded_rectangle((170, 142, 1110, 498), radius=32, outline=(10, 15, 27, 80), width=4)

    title_font = ImageFont.load_default(size=72)
    subtitle_font = ImageFont.load_default(size=38)

    draw.text((220, 220), "Indus Personal Work Track", fill=(15, 23, 42, 255), font=title_font)
    draw.text(
        (220, 320),
        "Offline-first project, issue and workflow manager",
        fill=(51, 65, 85, 255),
        font=subtitle_font,
    )
    draw.text((220, 380), "Go + Wails + React + SQLite", fill=(14, 116, 144, 255), font=subtitle_font)

    banner_path = DOC_ASSETS / "repository-banner.png"
    img.convert("RGB").save(banner_path, format="PNG", optimize=True)


def main() -> None:
    ensure_dirs()
    generate_app_icon()
    generate_repo_banner()
    print("Generated app icon and repository banner.")


if __name__ == "__main__":
    main()
