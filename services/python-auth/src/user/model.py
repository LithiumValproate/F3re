import bcrypt
from enum import Enum
from sqlalchemy import Column, String, Enum as SAEnum
from sqlalchemy.orm import declarative_base
import f3re


Base = declarative_base()


class UserType(str, Enum):
    STUDENT = "student"
    TEACHER = "teacher"
    ADMIN = "admin"
    BOT = "bot"


class User(Base):
    __tablename__ = 'users'

    id = Column(String(50), primary_key=True, index=True)
    name = Column(String(100), nullable=False)
    password_hash = Column(String, nullable=False)
    user_type = Column(SAEnum(UserType), nullable=False)

    def set_password(self, password: str):
        salt = bcrypt.gensalt()
        hashed_bytes = bcrypt.hashpw(password.encode('utf-8'), salt)
        self.password_hash = hashed_bytes.decode('utf-8')

    def check_password(self, password: str) -> bool:
        return bcrypt.checkpw(password.encode('utf-8'), self.password_hash.encode('utf-8'))

    def to_public_dict(self) -> dict:
        return {
            "id": self.id,
            "name": self.name,
            "user_type": self.user_type.value
        }
