import datetime as dt
import json
from enum import Enum
from typing import get_type_hints

from ..model.contact import Email, Phone
from ..model.models import Student
from . import utils


class EnhancedJSONDecoder(json.JSONDecoder):
    def __init__(self, *args, **kwargs):
        super().__init__(object_hook=self.object_hook, *args, **kwargs)

    def object_hook(self, dct):
        if '__type__' in dct:
            type_name = dct.pop('__type__')
            if type_name in utils.CLASS_REGISTRY:
                target_cls = utils.CLASS_REGISTRY[type_name]
                field_types = get_type_hints(target_cls)
                for key, value in dct.items():
                    if key in field_types:
                        field_type = field_types[key]
                        if field_type is Phone or field_type is Email:
                            dct[key] = field_type(value)
                        elif isinstance(field_type, type) and issubclass(field_type, Enum):
                            dct[key] = field_type(value)
                        elif field_type is dt.date and isinstance(value, str):
                            dct[key] = dt.date.fromisoformat(value)
                        elif field_type is dt.datetime and isinstance(value, str):
                            dct[key] = dt.datetime.fromisoformat(value)
                return target_cls(**dct)
        return dct


def student_from_json(json_string: str) -> Student:
    """
    Args:
        json_string: 包含学生数据的 JSON 格式字符串。
    Returns:
        一个 Student 类的实例。
    Raises:
        json.JSONDecodeError: 如果字符串不是有效的 JSON。
        ValueError: 如果解码后的对象不是一个 Student 实例。
    """
    decoded_object = json.loads(json_string, cls=EnhancedJSONDecoder)
    if not isinstance(decoded_object, Student):
        raise ValueError(
            f"JSON did not represent a Student object. "
            f"Decoded to type: {type(decoded_object).__name__}"
        )
    return decoded_object
