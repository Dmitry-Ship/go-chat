import React, { FormEvent } from "react";
import styles from "./Login.module.css";
import { Link } from "react-router-dom";
import { useAuth } from "../../authContext";

function Login() {
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const { login } = useAuth();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    login(username, password);
  };

  return (
    <>
      <header className={`header`}>
        <h2>Login</h2>
      </header>
      <section className={styles.wrap}>
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            name="username"
            className="input"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />

          <input
            type="password"
            name="password"
            className="input m-top-1"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <button
            disabled={!username || !password}
            type="submit"
            className="btn m-top-1"
          >
            Login
          </button>
          <Link to="/signup" className={`m-top-1 ${styles.signUpLink}`}>
            or Sign Up
          </Link>
        </form>
      </section>
    </>
  );
}

export default Login;
