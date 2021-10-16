import { useState, useEffect } from "react";
import { Register } from "./components/Register";
import { Login } from "./components/Login";

import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useParams,
} from "react-router-dom";

function App() {
  const [jwt, setJwt] = useState("");

  const handleJWTChange = (jwtToken: string) => {
    setJwt(jwtToken);
  };

  const logout = () => {
    setJwt("");
    window.localStorage.removeItem("jwt");
  };

  useEffect(() => {
    let t = window.localStorage.getItem("jwt");

    if (t) {
      if (jwt === "") {
        setJwt(JSON.parse(t));
      }
    }
  }, []);

  let loginLink;

  if (jwt === "") {
    loginLink = <Link to="/login">Login</Link>;
  } else {
    loginLink = (
      <Link to="logout" onClick={logout}>
        Logout
      </Link>
    );
  }

  return (
    <Router>
      <header className="App-header">
        <nav className="font-sans flex flex-col text-center sm:flex-row sm:text-left sm:justify-between py-4 px-6 bg-white shadow sm:items-baseline w-full">
          <div className="text-2xl no-underline text-grey-darkest hover:text-blue-dark">
            {loginLink}
          </div>
          <div className="mb-2 sm:mb-0">
            <Link
              to="/"
              className="text-2xl no-underline text-grey-darkest hover:text-blue-dark"
            >
              Home
            </Link>
          </div>
          <div>
            <Link
              to="/form/submission"
              className="text-lg no-underline text-grey-darkest hover:text-blue-dark ml-2"
            >
              Register
            </Link>
            <Link
              to="/login"
              className="text-lg no-underline text-grey-darkest hover:text-blue-dark ml-2"
            >
              Login
            </Link>
            <a
              href="/three"
              className="text-lg no-underline text-grey-darkest hover:text-blue-dark ml-2"
            >
              GitHub
            </a>
          </div>
        </nav>
      </header>
      <Switch>
        <Route path="/form/submission" component={Register}></Route>
        <Route
          path="/login"
          render={() => <Login jwtProps={handleJWTChange} />}
        ></Route>
      </Switch>
    </Router>
  );
}

export default App;
