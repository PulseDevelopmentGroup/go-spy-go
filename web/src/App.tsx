import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import styled from '@emotion/styled';

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

    this.state = {
      socket: null,
    };
  }

  componentDidMount() {
    const socket = new WebSocket(`ws://${process.env.REACT_APP_API_URL}/api`);

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
        console.log(e);
      };

      window.onbeforeunload = () => {
        console.log('firing');
        socket.close();
      };
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
