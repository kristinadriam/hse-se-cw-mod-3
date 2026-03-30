import subprocess


def test_help_exits_zero(chat_binary):
    r = subprocess.run(
        [str(chat_binary), "-h"],
        capture_output=True,
        text=True,
    )
    assert r.returncode == 0


def test_name_required(chat_binary):
    r = subprocess.run(
        [str(chat_binary)],
        capture_output=True,
        text=True,
    )
    assert r.returncode != 0
    assert "name" in (r.stderr + r.stdout).lower()


def test_invalid_connect_address(chat_binary):
    r = subprocess.run(
        [str(chat_binary), "-name", "x", "-connect", "not-a-host:port"],
        capture_output=True,
        text=True,
    )
    assert r.returncode != 0
