import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import styled from '@emotion/styled';
import SocketManager, { Action } from './store/SocketManager';

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

interface Props {}

export default class App extends Component<Props, any> {
  constructor(props: Props) {
    super(props);

    this.state = {};
  }

  componentDidMount() {
    const socket = new WebSocket(`ws://${process.env.REACT_APP_API_URL}/api`);

    socket.onopen = e => {
      SocketManager.init(socket);

      let message = {
        kind: 'CREATE_GAME',
        data: '{"code":"", "username":"user"}',
      } as Action<any>;

      SocketManager.sendMessage(message);
    };
  }

  render() {
    return (
      <Router>
        <Container>
          <Route exact path="/" component={Landing} />
        </Container>
      </Router>
    );
  }
}
