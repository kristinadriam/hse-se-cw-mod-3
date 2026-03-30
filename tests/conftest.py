import subprocess
from pathlib import Path

import pytest


@pytest.fixture(scope="session")
def project_root() -> Path:
    return Path(__file__).resolve().parent.parent


@pytest.fixture(scope="session")
def chat_binary(project_root: Path) -> Path:
    bin_dir = project_root / "bin"
    bin_dir.mkdir(exist_ok=True)
    out = bin_dir / "chat"
    subprocess.run(
        ["go", "build", "-o", str(out), "./cmd/chat"],
        cwd=project_root,
        check=True,
    )
    return out
