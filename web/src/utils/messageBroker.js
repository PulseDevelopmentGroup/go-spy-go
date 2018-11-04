class MessageBroker {
  handleMessage(message) {
    const parsedMessage = this.parseMessage(message);
    console.log(parsedMessage);
  }
  
  parseMessage(e) {
    const parsedMessage = JSON.parse(e.data);
    try {
      return {
        type: parsedMessage.kind,
        data: JSON.parse(parsedMessage.data)
      }
    } catch (err) {
      console.log('Error occured while parsing socket message:', err);
    }
  }
  
}

export default new MessageBroker();