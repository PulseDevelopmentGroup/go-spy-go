import React, { Component } from 'react';
import styled from '@emotion/styled';

import { Button } from './components/Button';
import { GameInfo } from './models';

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

const StyledTextField = styled.input`
  margin: 10px 15px;
  padding: 5px 10px;

  border: 2px solid #4392f1;
  font-family: 'Oswald', sans-serif;
  font-size: 1.2em;
`;

interface ControlsProps {
  gameId: string;
  userName: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const CreateControls: React.FC<ControlsProps> = props => {
  return (
    <>
      <StyledTextField
        placeholder="Game ID"
        name="gameId"
        type="text"
        value={props.gameId}
        onChange={props.onChange}
      />
      <StyledTextField
        placeholder="Username"
        name="userName"
        type="text"
        value={props.userName}
        onChange={props.onChange}
      />
    </>
  );
};

const JoinControls: React.FC<ControlsProps> = props => {
  return (
    <>
      <StyledTextField
        placeholder="Game ID"
        name="gameId"
        type="text"
        value={props.gameId}
        onChange={props.onChange}
      />
      <StyledTextField
        placeholder="Username"
        name="userName"
        type="text"
        value={props.userName}
        onChange={props.onChange}
      />
    </>
  );
};

const SubmitBtn = styled(Button)``;

interface Props {}

interface State {
  gameInfo: GameInfo;
  showCreate: boolean;
  showJoin: boolean;
}

export default class Landing extends Component<Props, State> {
  form: React.RefObject<HTMLFormElement>;
  constructor(props: Props) {
    super(props);

    this.state = {
      showCreate: false,
      showJoin: false,
      gameInfo: {
        gameId: '',
        userName: '',
      },
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

  changeHandler(e: React.ChangeEvent<HTMLInputElement>) {
    const name = e.target.name;
    const value = e.target.value;

    let newGameInfo = {
      ...this.state.gameInfo,
    };

    newGameInfo[name as keyof GameInfo] = value;

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
          <Button
            className={this.state.showCreate ? 'active' : ''}
            onClick={this.showCreateFields}
          >
            Create
          </Button>
          <Button
            className={this.state.showJoin ? 'active' : ''}
            onClick={this.showJoinFields}
          >
            Join
          </Button>
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
            <SubmitBtn type="submit">Submit</SubmitBtn>
          )}
        </ConnectFields>
      </StyledLanding>
    );
  }
}
