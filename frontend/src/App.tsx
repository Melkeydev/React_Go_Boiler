import { useState } from "react";
import { EditData } from "./components/EditData";

import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useParams,
} from "react-router-dom";

function App() {
  return (
    <Router>
      <header className="App-header">
        <nav className="font-sans flex flex-col text-center sm:flex-row sm:text-left sm:justify-between py-4 px-6 bg-white shadow sm:items-baseline w-full">
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
              Form
            </Link>
            <a
              href="/two"
              className="text-lg no-underline text-grey-darkest hover:text-blue-dark ml-2"
            >
              About
            </a>
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
        <Route path="/form/submission" component={EditData}></Route>
      </Switch>
    </Router>
  );
}

export default App;
