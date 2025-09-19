from dataclasses import dataclass
import re


@dataclass(frozen=True)
class Phone:
    value: str

    def __post_init__(self):
        if not isinstance(self.value, str):
            raise TypeError('Initial value must be a string.')
        # A simple check for Chinese mobile phone numbers
        if not re.match(r"^\d{11}$", self.value):
            raise ValueError("Invalid phone number format. Expected 11 digits.")


@dataclass(frozen=True)
class Email:
    value: str

    def __post_init__(self):
        if not isinstance(self.value, str):
            raise TypeError('Initial value must be a string.')
        if not re.match(r"^[\w.-]+@[\w.-]+\.\w+$", self.value):
            raise ValueError('Invalid email format.')
