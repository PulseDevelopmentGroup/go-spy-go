/**
 *  Create a payload object for transmission over socket
 * @param {string} kind - The type of message
 * @param {string} data - The payload of the message
 */
function packMessage(kind, data) {
  return JSON.stringify({
    kind,
    data
  });
}

export {
  packMessage
}