import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import styled from 'styled-components';

import { packMessage } from './utils/socketUtils';
import MessageBroker from './utils/messageBroker';

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

    this.state = {};

    this.socket = null;

    this.createGame = this.createGame.bind(this);
    this.joinGame = this.joinGame.bind(this);

    this.gameFunctions = {
      createGame: this.createGame,
      joinGame: this.joinGame,
    };
  }

  componentDidMount() {
    const socket = new WebSocket(`ws://${process.env.API_URL}/api`);

    socket.onopen = e => {
      this.socket = socket;

      let message = {
        type: 'create-game',
        data: '{"code":"", "username":"user"}',
      };

      let obj = JSON.stringify(message);

      socket.onerror = err => {
        console.log(`Following error occured with websocket: ${err}`);
      };

      socket.onmessage = e => {
        MessageBroker.handleMessage(e);
      };
      socket.send(obj);

      window.onbeforeunload = () => {
        console.log('firing');
        socket.close();
      };
    };
  }

  onSocketMessage(message) {
    console.log('Message recieved');
    console.log(message);
  }

  createGame(id, username) {
    const gameObj = {
      gameId: id || '',
      username: username,
    };

    const payload = packMessage('CREATE_GAME', JSON.stringify(gameObj));
    console.log(payload);
    this.socket.send(payload);
  }

  joinGame(id, username) {
    const gameObj = {
      gameId: id || '',
      username: username,
    };

    const payload = packMessage('JOIN_GAME', JSON.stringify(gameObj));
    console.log(payload);
    this.socket.send(payload);
  }

  render() {
    return (
      <Router>
        <Container>
          <GlobalStyles />
          <Route
            exact
            path="/"
            render={props => (
              <Landing {...props} gameFunctions={this.gameFunctions} />
            )}
          />
        </Container>
      </Router>
    );
  }
}
