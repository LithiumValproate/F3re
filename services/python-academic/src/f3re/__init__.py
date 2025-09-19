"""Namespace package for the F3re project."""

from importlib import import_module
from importlib.util import find_spec
from pkgutil import extend_path
from typing import Any

__path__ = extend_path(__path__, __name__)
__all__: list[str] = []


def __getattr__(name: str) -> Any:
    qualified_name = f"{__name__}.{name}"
    if find_spec(qualified_name) is None:
        raise AttributeError(f"module '{__name__}' has no attribute '{name}'")

    module = import_module(qualified_name)
    globals()[name] = module
    if name not in __all__:
        __all__.append(name)
    return module
