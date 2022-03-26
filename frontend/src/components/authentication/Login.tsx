import React, { FormEvent } from "react";
import styles from "./Login.module.css";
import Link from "next/link";
import { useAuth } from "../../contexts/authContext";

function Login() {
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const { login } = useAuth();

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    login(username, password);
  };

  return (
    <>
      <header className={`header`}>
        <h2>Go Chat</h2>
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

          <Link href="/signup">
            <a className={`m-top-1 ${styles.signUpLink}`}>or Sign Up</a>
          </Link>
        </form>
      </section>
    </>
  );
}

export default Login;
