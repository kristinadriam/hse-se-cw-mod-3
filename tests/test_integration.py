import os
import select
import signal
import socket
import subprocess
import time


def _free_port() -> int:
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(("127.0.0.1", 0))
    port = s.getsockname()[1]
    s.close()
    return port


def _wait_stderr_contains(proc: subprocess.Popen, needle: str, timeout: float) -> None:
    buf = b""
    deadline = time.time() + timeout
    while time.time() < deadline:
        if needle.encode() in buf or needle in buf.decode(errors="replace"):
            return
        remaining = deadline - time.time()
        if remaining <= 0:
            break
        r, _, _ = select.select([proc.stderr], [], [], min(0.2, remaining))
        if r:
            buf += os.read(proc.stderr.fileno(), 4096)
    raise AssertionError(
        f"stderr did not contain {needle!r} in time; got: {buf.decode(errors='replace')!r}"
    )


def _read_stdout_until(
    proc: subprocess.Popen, needle: str, timeout: float
) -> str:
    buf = b""
    deadline = time.time() + timeout
    while time.time() < deadline:
        if needle in buf.decode(errors="replace"):
            return buf.decode(errors="replace")
        remaining = deadline - time.time()
        if remaining <= 0:
            break
        r, _, _ = select.select([proc.stdout], [], [], min(0.2, remaining))
        if r:
            buf += os.read(proc.stdout.fileno(), 4096)
    raise AssertionError(
        f"stdout did not contain {needle!r} in time; got: {buf.decode(errors='replace')!r}"
    )


def _terminate(proc: subprocess.Popen) -> None:
    if proc.poll() is not None:
        return
    proc.send_signal(signal.SIGTERM)
    try:
        proc.wait(timeout=5)
    except subprocess.TimeoutExpired:
        proc.kill()
        proc.wait(timeout=5)


def test_listener_and_client_exchange_message(chat_binary):
    port = _free_port()
    addr = f"127.0.0.1:{port}"

    alice = subprocess.Popen(
        [
            str(chat_binary),
            "-name",
            "Alice",
            "-listen",
            addr,
        ],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    try:
        _wait_stderr_contains(alice, "listening", timeout=15)

        bob = subprocess.Popen(
            [
                str(chat_binary),
                "-name",
                "Bob",
                "-connect",
                addr,
            ],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        try:
            time.sleep(0.3)
            assert alice.stdin is not None
            alice.stdin.write(b"Hello from Alice\n")
            alice.stdin.flush()

            out = _read_stdout_until(bob, "Hello from Alice", timeout=15)
            assert "Alice" in out
        finally:
            _terminate(bob)
    finally:
        _terminate(alice)
