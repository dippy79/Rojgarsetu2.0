import React from 'react';
import { Link, NavLink } from 'react-router-dom';
import './Navbar.css';

const Navbar = () => {
  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <Link to="/">RojgarSetu</Link>
      </div>
      <ul className="navbar-links">
        <li>
          <NavLink to="/" end>Home</NavLink>
        </li>
        <li>
          <NavLink to="/gov-jobs">Gov Jobs</NavLink>
        </li>
        <li>
          <NavLink to="/private-jobs">Private Jobs</NavLink>
        </li>
        <li>
          <NavLink to="/courses">Courses</NavLink>
        </li>
        <li>
          <NavLink to="/videos">Videos</NavLink>
        </li>
      </ul>
    </nav>
  );
};

export default Navbar;

