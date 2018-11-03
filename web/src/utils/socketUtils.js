/**
 *  Create a payload object for transmission over socket
 * @param {string} trigger - The type of message
 * @param {string} data - The payload of the message
 */
function packMessage(trigger, data) {
  return JSON.stringify({
    trigger,
    data
  });
}

export {
  packMessage
}