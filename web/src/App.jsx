import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import styled from 'styled-components';

import GlobalStyles from './components/Global';
import Landing from './Landing';

const Container = styled.div`
  display: flex;
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;

  flex-direction: column;
  justify-content: center;
  padding: 15px;

  background-color: #fdfdff;
`;

export default class App extends Component {
  constructor() {
    super();

    this.state = {
      socket: null,
    };
  }

  componentDidMount() {
    const socket = new WebSocket(`ws://${process.env.API_URL}/api`);

    socket.onopen = e => {
      this.setState({
        socket,
      });

      let message = {
        type: 'create-game',
        data: '{"code":"", "username":"user"}',
      };

      let obj = JSON.stringify(message);

      socket.send(obj);
      socket.onmessage = e => {
        console.log(e.data);
      };
    };
  }

  render() {
    return (
      <Router>
        <Container>
          <GlobalStyles />
          <Route exact path="/" render={props => <Landing {...props} />} />
        </Container>
      </Router>
    );
  }
}
