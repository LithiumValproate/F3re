from fastapi import FastAPI

app = FastAPI(
    title="F3re Authentication Service",
    description="Handles user creation, authentication, and authorization.",
    version="0.1.0"
)

@app.get("/")
def read_root():
    return {"message": "F3re Authentication Service is running"}

# In a real application, you would add endpoints for user registration, login, etc.
# For example:
#
# from .user import User
#
# @app.post("/users/")
# def create_user(user: User):
#     # Logic to create a user
#     return {"message": f"User {user.name} created"}
