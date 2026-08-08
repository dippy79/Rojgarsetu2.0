import React from 'react';
import { Link, NavLink } from 'react-router-dom';
import './Navbar.css';

// Flat, high-visibility navbar. No dropdowns — every section is a visible,
// side-by-side link for maximum discoverability (desktop + mobile).
const Navbar = () => {
  return (
    <nav className="navbar" aria-label="Primary navigation">
      <div className="navbar-brand">
        <Link to="/">🏛️ RojgarSetu</Link>
      </div>
      <ul className="navbar-links">
        <li>
          <NavLink to="/" end>🏠 Home</NavLink>
        </li>
        <li>
          <NavLink to="/gov-jobs">🏛️ Govt Jobs</NavLink>
        </li>
        <li>
          <NavLink to="/private-jobs">🏢 Private Jobs</NavLink>
        </li>
        <li>
          <NavLink to="/courses">📚 Courses</NavLink>
        </li>
        <li>
          <NavLink to="/videos">🎥 Videos</NavLink>
        </li>
        <li>
          <NavLink to="/govt-forms">🗂️ Govt Forms</NavLink>
        </li>
      </ul>
    </nav>
  );
};

export default Navbar;
