from collections import UserString
from dataclasses import dataclass
import re


class Phone(UserString):
    def __init__(self, value: str, max_length: int = 11):
        if not isinstance(value, str):
            raise TypeError('Initial value must be a string.')
        if len(value) > max_length:
            raise ValueError(f"Initial string exceeds the maximum length of {max_length}.")
        super().__init__(value)


@dataclass
class Email(UserString):
    value: str

    def __post_init__(self):
        if not isinstance(self.value, str):
            raise TypeError('Initial value must be a string.')
        if not re.match(r"^[\w.-]+@[\w.-]+\.\w+$", self.value):
            raise ValueError('Invalid email format.')
        super().__init__(self.value)
