/**
 * Anticipates the next block of data to read from serial.
 */
enum NextInputStep {
  MessageType,
  Target,
  Command,
  Data
};

boolean testMode = false;        // whether we're in automated testing mode
boolean messageReady = false;
int nextStep = MessageType;
String processed = "";
int messages = 0;

/**
 * The parts of the incoming message.
 */
String type = "";
String target = "";
String command = "";
String data = "";

/**
 * Send the command data to the appropriate instrument.
 */
void processMessage() {
  processed += type + ":" + target + ":" + command + "=" + data + ",";
  if (type.equals("SHIP")) {
    // Handle ship-specific inputs
    if (target.equals("FOCUS")) {
      // Message is for the currently-focused vessel
    } else {
      // Message is for another vessel
    }
  } else if (type.equals("CAMERA")) {
    // Handle camera-related messages
  } else if (type.equals("NAV")) {
    // Handle nav-related messages
  } else {
    // Handle all other types
  }
}

void setup() {
  // initialize serial
  Serial.begin(9600);
  // reserve bytes for the message parts
  type.reserve(30);
  target.reserve(30);
  command.reserve(100);
  data.reserve(200);
  processed.reserve(500);
}

void loop() {
  if (testMode) {
    testLoop();
  } else if (messageReady) {
    processMessage();
    Serial.println(processed);
    Serial.println(messages);
    resetState();
  }
}

void resetState() {
  nextStep = MessageType;
  type = "";
  target = "";
  command = "";
  data = "";
  messageReady = false;
}

/**
 SerialEvent occurs whenever a new data comes in the
 hardware serial RX.  This routine is run between each
 time loop() runs, so using delay inside loop can delay
 response.  Multiple bytes of data may be available.
 */
void serialEvent() {
  while (!messageReady && Serial.available()) {
    char c = (char)Serial.read();

    if (c == '\n') {
      // A forced way to re-orient in case of command misalignment
      processed = "";
      resetState();
    } else if (c == ':' || c == '=' || c == '|') {
      // Message part boundary
      nextStep++;
      if (type.equals("ORB") || type.equals("CAMERA") || type.equals("OBJ")) {
        // These types do not have a target, only a command.
        nextStep++;
      }
      // The pipe means we have a completed message
      if (c == '|') {
        messageReady = true;
        messages++;
      }
    } else {
      switch(nextStep) {
        case MessageType:
          // The first block identifies the type of message coming in.
          type += c;
          break;
        case Target:
          // Some message types have a target.
          target += c;
          break;
        case Command:
          // The command related to the incoming message.
          command += c;
          break;
        case Data:
          data += c;
          break;
      }
    }
  }
}

/******************************************************************************/
/********************           TEST CODE          ****************************/
/******************************************************************************/

struct TestState {
  boolean waitingForReply = false;
  String expectedReply = "";
  unsigned long lastHeartbeat = millis(); // the last time a command was issued
};

TestState testState;
unsigned long PULSE_FREQ = 2000; // every two seconds

/**
 * Toggles test mode, with confirmation back to the requester.
 */
void toggleTestMode() {
  testMode ^= true;
  // Send confirmation back
  Serial.println(testMode, DEC);
}

/**
 * Executed on each loop iteration when in automated test mode.
 */
void testLoop() {
  /*
  if (testState.waitingForReply && stringComplete) {
    inputString.trim();
    Serial.println(inputString.equals(testState.expectedReply));
    Serial.print("Lag ");
    Serial.println(millis() - testState.lastHeartbeat);
    testState.waitingForReply = false;
    processString();
  }
  if (!testState.waitingForReply && millis() >= testState.lastHeartbeat + PULSE_FREQ) {
    sendTestCommand();
  }
  */
}

void sendTestCommand() {
  String rndstr = "";
  for (int i=0; i < 10; i++) {
    byte randomValue = random(0, 37);
    char letter = randomValue + 'a';
    if (randomValue > 26) {
      letter = (randomValue - 26) + '0';
    }
    rndstr += letter;
  }
  testState.expectedReply = rndstr;

  Serial.print("TEST ");
  Serial.println(testState.expectedReply);
  testState.waitingForReply = true;
  testState.lastHeartbeat = millis();
}


