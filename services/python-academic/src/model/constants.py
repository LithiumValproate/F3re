from enum import Enum


class DayOfWeek(Enum):
    MONDAY = 1
    TUESDAY = 2
    WEDNESDAY = 3
    THURSDAY = 4
    FRIDAY = 5
    SATURDAY = 6
    SUNDAY = 7


class Repetition(Enum):
    WEEKLY = 1
    BIWEEKLY_ODD = 2
    BIWEEKLY_EVEN = 3


class Sex(Enum):
    MALE = 1
    FEMALE = 2
    NON_BINARY = 3


class Status(Enum):
    ACTIVE = 1
    INACTIVE = 2
    GRADUATED = 3


CHINA_PROVINCES = {
    "Beijing", "Tianjin", "Shanghai", "Chongqing",
    "Hebei", "Shanxi", "Liaoning", "Jilin", "Heilongjiang",
    "Jiangsu", "Zhejiang", "Anhui", "Fujian", "Jiangxi",
    "Shandong", "Henan", "Hubei", "Hunan", "Guangdong",
    "Guangxi", "Hainan", "Sichuan", "Guizhou", "Yunnan",
    "Tibet", "Shaanxi", "Gansu", "Qinghai", "Ningxia",
    "Xinjiang", "Inner Mongolia", "Hong Kong", "Macau", "Taiwan"
}
