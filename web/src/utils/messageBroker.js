class MessageBroker {
  handleMessage(message) {
    const parsedMessage = this.parseMessage(message);
    console.log(parsedMessage);
  }
  
  parseMessage(e) {
    debugger;
    const parsedMessage = JSON.parse(e.data);
    return {
      type: parsedMessage.Kind,
      data:parsedMessage.Data
    }
  }
  
}

export default new MessageBroker();