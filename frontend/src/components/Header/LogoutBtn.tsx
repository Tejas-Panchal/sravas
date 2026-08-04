import { useNavigate } from "react-router-dom";
import { Button } from "../index.ts";
import { useAppDispatch } from "../../store/hooks.ts";
import { logout } from "../../store/authSlice.ts";
import authService from "../../api/auth.ts";

function LogoutBtn() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  const logoutHandler = async () => {
    await authService.logout();
    dispatch(logout());
    navigate("/");
  };

  return (
    <Button bgColor="bg-gray-800" onClick={logoutHandler}>
      Logout
    </Button>
  );
}

export default LogoutBtn;
