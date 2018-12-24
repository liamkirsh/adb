<template>
  <div>
    <div id="flash" role="alert"></div>
    <nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container">
        <div class="navbar-header">
          <button
            type="button"
            class="navbar-toggle collapsed"
            data-toggle="collapse"
            data-target="#navbar"
            aria-expanded="false"
            aria-controls="navbar"
          >
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <span
            class="navbar-brand"
            href="#"
            style="color: white; font-weight: bold;"
          >Activist Database</span>
        </div>
        <div id="navbar" class="collapse navbar-collapse">
          <ul class="nav navbar-nav">
            <li class="dropdown">
              <a class="dropdown-toggle" data-toggle="dropdown" href="#">
                Events
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu">
                <li>
                  <router-link to="/">New Event</router-link>
                </li>
                <li>
                  <router-link to="/list_events">All Events</router-link>
                </li>
              </ul>
            </li>
            <li class="dropdown">
              <a class="dropdown-toggle" data-toggle="dropdown" href="#">
                Connections
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu">
                <li>
                  <router-link to="/new_connection">New Maintenance Connection</router-link>
                </li>
                <li>
                  <router-link to="/list_connections">All Maintenance Connections</router-link>
                </li>
                <li>
                  <router-link to="/activist_pool">Recruitment Connections</router-link>
                </li>
              </ul>
            </li>
            <li class="dropdown">
              <a class="dropdown-toggle" data-toggle="dropdown" href="#">
                Circles
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu">
                <li>
                  <router-link to="/circle_member_prospects">Circle Member Prospects</router-link>
                </li>
                <li>
                  <router-link to="/list_circles">Circles</router-link>
                </li>
              </ul>
            </li>
            <li class="dropdown">
              <a class="dropdown-toggle" data-toggle="dropdown" href="#">
                Chapter Members
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu">
                <li>
                  <router-link to="/chapter_member_prospects">Chapter Member Prospects</router-link>
                </li>
                <li>
                  <router-link to="/chapter_member_development">Chapter Member Development</router-link>
                </li>
              </ul>
            </li>
            <li class="dropdown">
              <a class="dropdown-toggle" data-toggle="dropdown" href="#">
                Organizers
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu">
                <li>
                  <router-link to="/organizer_prospects">Organizer Prospects</router-link>
                </li>
                <li>
                  <router-link to="/activist_development">Organizer Development</router-link>
                </li>
                <li>
                  <router-link to="/list_working_groups">Working Groups</router-link>
                </li>
              </ul>
            </li>
            <li class="dropdown">
              <a class="dropdown-toggle" data-toggle="dropdown" href="#">
                All
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu">
                <li>
                  <router-link to="/list_activists">All Activists</router-link>
                </li>
                <li>
                  <router-link to="/leaderboard">Leaderboard</router-link>
                </li>
                <!---<li><router-link to="/activist_recruitment">Recruitment</router-link></li>-->
                <!--<li><router-link to="/activist_actionteam">Action Team</router-link></li>-->
              </ul>
            </li>

            <!--<li><router-link to="/power">MPI</router-link></li>-->
            <!-- TODO(mdempsky): Hide for non-admins. -->
            <li>
              <router-link to="/admin/users">Users</router-link>
            </li>
          </ul>
        </div>
        <!--/.nav-collapse -->
      </div>
    </nav>

    <router-view></router-view>

    <center>
      <br>
      <img id="footer-img" src="/static/img/logo1.png" alt="DxE" width="50" height="50">
      <br>
      <br>
    </center>
  </div>
</template>

<script>
import { initializeFlashMessage } from "flash_message.js";

initializeFlashMessage();

export default {};
</script>
