from ..model import models as ac, contact

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