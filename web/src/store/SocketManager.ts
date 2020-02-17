export interface Action<T> {
  kind: string;
  data: T;
}

class SocketManager {
  private ws?: WebSocket;
  init(socket: WebSocket) {
    this.ws = socket;

    this.ws.onmessage = this.handleMessage;

    window.addEventListener('beforeunload', this.teardown);
  }

  handleMessage(e: MessageEvent) {
    console.log('got message');
    console.log(e);
  }

  sendMessage(message: Action<any>) {
    if (!this.ws) {
      throw new Error('Unable to send socket message before initialization');
    }

    let payload;
    try {
      payload = JSON.stringify(message);
    } catch (e) {
      console.error('Unable to serialize message payload', e);
      return;
    }

    this.ws.send(payload);
  }

  teardown() {
    this.ws?.close();
  }
}

export default new SocketManager();
