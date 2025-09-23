from f3re.academic import models as ac
from f3re.academic.model import contact

CLASS_REGISTRY = {
    'Student': ac.Student,
    'Course': ac.Course,
    'Teacher': ac.Teacher,
    'Address': ac.Address,
    'FamilyMember': ac.FamilyMember,
    'Grade': ac.Grade,
    'TimeSlot': ac.TimeSlot,
    'Phone': contact.Phone,
    'Email': contact.Email,
}