import React, { Component } from 'react';
import styled from 'styled-components';

import { Button } from './components/Button';

const StyledLanding = styled.div``;

const LandingHeader = styled.h1`
  margin-bottom: 0;
  text-align: center;

  color: #212922;
  font-family: 'Oswald', sans-serif;
  font-size: 4em;
`;

const Subheader = styled.h3`
  margin-top: 5px;
  text-align: center;

  color: #c1c1c1;
  font-family: 'Oswald', sans-serif;
  font-size: 1.5em;
`;

const ButtonGroup = styled.div`
  display: flex;
  justify-content: center;
`;

const ConnectFields = styled.form`
  display: flex;
  flex-direction: column;

  align-items: center;
`;

const StyledTextField = styled.input.attrs({
  type: 'text',
})`
  margin: 10px 15px;
  padding: 5px 10px;

  border: 2px solid #4392f1;
  font-family: 'Oswald', sans-serif;
  font-size: 1.2em;
`;

const CreateControls = props => {
  return (
    <>
      <StyledTextField
        placeholder="Game ID"
        name="gameId"
        value={props.gameId}
        onChange={props.onChange}
      />
      <StyledTextField
        placeholder="Username"
        name="userName"
        value={props.userName}
        onChange={props.onChange}
      />
    </>
  );
};

const JoinControls = props => {
  return (
    <>
      {console.log(props)}
      <StyledTextField
        placeholder="Game ID"
        name="gameId"
        value={props.gameId}
        onChange={props.onChange}
      />
      <StyledTextField
        placeholder="Username"
        name="userName"
        value={props.userName}
        onChange={props.onChange}
      />
    </>
  );
};

const SubmitBtn = styled(Button)``;

export default class Landing extends Component {
  constructor() {
    super();

    this.state = {
      showCreate: false,
      showJoin: false,
      gameInfo: {},
    };

    this.form = React.createRef();

    this.showCreateFields = this.showCreateFields.bind(this);
    this.showJoinFields = this.showJoinFields.bind(this);
    this.sendGameRequest = this.sendGameRequest.bind(this);
    this.changeHandler = this.changeHandler.bind(this);
  }

  showCreateFields() {
    this.setState({
      showCreate: true,
      showJoin: false,
    });
  }

  showJoinFields() {
    this.setState({
      showCreate: false,
      showJoin: true,
    });
  }

  sendGameRequest() {
    console.log(this.form.current);
    debugger;
  }

  changeHandler(e) {
    const name = e.target.name;
    const value = e.target.value;

    let newGameInfo = this.state.gameInfo;

    newGameInfo[name] = value;

    this.setState({
      gameInfo: newGameInfo,
    });
  }

  render() {
    return (
      <StyledLanding>
        <LandingHeader>Welcome to Go Spy Go!</LandingHeader>
        <Subheader>
          An implementation of{' '}
          <a href="http://international.hobbyworld.ru/spyfall">Spyfall</a>
        </Subheader>
        <ButtonGroup>
          <Button onClick={this.showCreateFields}>Create</Button>
          <Button onClick={this.showJoinFields}>Join</Button>
        </ButtonGroup>
        <ConnectFields
          onSubmit={e => {
            e.preventDefault();
            this.sendGameRequest();
          }}
          ref={this.form}
        >
          {this.state.showCreate && (
            <CreateControls
              gameId={
                this.state.gameInfo.gameId ? this.state.gameInfo.gameId : ''
              }
              userName={
                this.state.gameInfo.userName ? this.state.gameInfo.userName : ''
              }
              onChange={this.changeHandler}
            />
          )}
          {this.state.showJoin && (
            <JoinControls
              gameId={
                this.state.gameInfo.gameId ? this.state.gameInfo.gameId : ''
              }
              userName={
                this.state.gameInfo.userName ? this.state.gameInfo.userName : ''
              }
              onChange={this.changeHandler}
            />
          )}
          {(this.state.showCreate || this.state.showJoin) && (
            <SubmitBtn as="input" type="submit" value="Submit" />
          )}
        </ConnectFields>
      </StyledLanding>
    );
  }
}
