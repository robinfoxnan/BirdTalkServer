// source: user.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

goog.provide('proto.model.FriendOpResult');

goog.require('jspb.BinaryReader');
goog.require('jspb.BinaryWriter');
goog.require('jspb.Map');
goog.require('jspb.Message');
goog.require('proto.model.UserInfo');

goog.forwardDeclare('proto.model.UserOperationType');
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.model.FriendOpResult = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.model.FriendOpResult.repeatedFields_, null);
};
goog.inherits(proto.model.FriendOpResult, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.FriendOpResult.displayName = 'proto.model.FriendOpResult';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.model.FriendOpResult.repeatedFields_ = [4];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.model.FriendOpResult.prototype.toObject = function(opt_includeInstance) {
  return proto.model.FriendOpResult.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.FriendOpResult} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.FriendOpResult.toObject = function(includeInstance, msg) {
  var f, obj = {
    operation: jspb.Message.getFieldWithDefault(msg, 1, 0),
    result: jspb.Message.getFieldWithDefault(msg, 2, ""),
    user: (f = msg.getUser()) && proto.model.UserInfo.toObject(includeInstance, f),
    usersList: jspb.Message.toObjectList(msg.getUsersList(),
    proto.model.UserInfo.toObject, includeInstance),
    sendid: jspb.Message.getFieldWithDefault(msg, 5, 0),
    msgid: jspb.Message.getFieldWithDefault(msg, 6, 0),
    paramsMap: (f = msg.getParamsMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.model.FriendOpResult}
 */
proto.model.FriendOpResult.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.FriendOpResult;
  return proto.model.FriendOpResult.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.FriendOpResult} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.FriendOpResult}
 */
proto.model.FriendOpResult.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.model.UserOperationType} */ (reader.readEnum());
      msg.setOperation(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setResult(value);
      break;
    case 3:
      var value = new proto.model.UserInfo;
      reader.readMessage(value,proto.model.UserInfo.deserializeBinaryFromReader);
      msg.setUser(value);
      break;
    case 4:
      var value = new proto.model.UserInfo;
      reader.readMessage(value,proto.model.UserInfo.deserializeBinaryFromReader);
      msg.addUsers(value);
      break;
    case 5:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setSendid(value);
      break;
    case 6:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setMsgid(value);
      break;
    case 7:
      var value = msg.getParamsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.model.FriendOpResult.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.FriendOpResult.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.FriendOpResult} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.FriendOpResult.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOperation();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getResult();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getUser();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.model.UserInfo.serializeBinaryToWriter
    );
  }
  f = message.getUsersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.model.UserInfo.serializeBinaryToWriter
    );
  }
  f = message.getSendid();
  if (f !== 0) {
    writer.writeInt64(
      5,
      f
    );
  }
  f = message.getMsgid();
  if (f !== 0) {
    writer.writeInt64(
      6,
      f
    );
  }
  f = message.getParamsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(7, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional UserOperationType operation = 1;
 * @return {!proto.model.UserOperationType}
 */
proto.model.FriendOpResult.prototype.getOperation = function() {
  return /** @type {!proto.model.UserOperationType} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.model.UserOperationType} value
 * @return {!proto.model.FriendOpResult} returns this
 */
proto.model.FriendOpResult.prototype.setOperation = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional string result = 2;
 * @return {string}
 */
proto.model.FriendOpResult.prototype.getResult = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.model.FriendOpResult} returns this
 */
proto.model.FriendOpResult.prototype.setResult = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional UserInfo user = 3;
 * @return {?proto.model.UserInfo}
 */
proto.model.FriendOpResult.prototype.getUser = function() {
  return /** @type{?proto.model.UserInfo} */ (
    jspb.Message.getWrapperField(this, proto.model.UserInfo, 3));
};


/**
 * @param {?proto.model.UserInfo|undefined} value
 * @return {!proto.model.FriendOpResult} returns this
*/
proto.model.FriendOpResult.prototype.setUser = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.model.FriendOpResult} returns this
 */
proto.model.FriendOpResult.prototype.clearUser = function() {
  return this.setUser(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.model.FriendOpResult.prototype.hasUser = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * repeated UserInfo users = 4;
 * @return {!Array<!proto.model.UserInfo>}
 */
proto.model.FriendOpResult.prototype.getUsersList = function() {
  return /** @type{!Array<!proto.model.UserInfo>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.UserInfo, 4));
};


/**
 * @param {!Array<!proto.model.UserInfo>} value
 * @return {!proto.model.FriendOpResult} returns this
*/
proto.model.FriendOpResult.prototype.setUsersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.model.UserInfo=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.UserInfo}
 */
proto.model.FriendOpResult.prototype.addUsers = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.model.UserInfo, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.FriendOpResult} returns this
 */
proto.model.FriendOpResult.prototype.clearUsersList = function() {
  return this.setUsersList([]);
};


/**
 * optional int64 sendId = 5;
 * @return {number}
 */
proto.model.FriendOpResult.prototype.getSendid = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 5, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.FriendOpResult} returns this
 */
proto.model.FriendOpResult.prototype.setSendid = function(value) {
  return jspb.Message.setProto3IntField(this, 5, value);
};


/**
 * optional int64 msgId = 6;
 * @return {number}
 */
proto.model.FriendOpResult.prototype.getMsgid = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.FriendOpResult} returns this
 */
proto.model.FriendOpResult.prototype.setMsgid = function(value) {
  return jspb.Message.setProto3IntField(this, 6, value);
};


/**
 * map<string, string> params = 7;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.model.FriendOpResult.prototype.getParamsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 7, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.model.FriendOpResult} returns this
 */
proto.model.FriendOpResult.prototype.clearParamsMap = function() {
  this.getParamsMap().clear();
  return this;};


