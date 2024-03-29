// source: msg.proto
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

goog.provide('proto.model.MsgChatQueryResult');

goog.require('jspb.BinaryReader');
goog.require('jspb.BinaryWriter');
goog.require('jspb.Message');
goog.require('proto.model.MsgChat');
goog.require('proto.model.MsgChatReply');

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
proto.model.MsgChatQueryResult = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.model.MsgChatQueryResult.repeatedFields_, null);
};
goog.inherits(proto.model.MsgChatQueryResult, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.model.MsgChatQueryResult.displayName = 'proto.model.MsgChatQueryResult';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.model.MsgChatQueryResult.repeatedFields_ = [3,4];



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
proto.model.MsgChatQueryResult.prototype.toObject = function(opt_includeInstance) {
  return proto.model.MsgChatQueryResult.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.model.MsgChatQueryResult} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.MsgChatQueryResult.toObject = function(includeInstance, msg) {
  var f, obj = {
    userid: jspb.Message.getFieldWithDefault(msg, 1, 0),
    toid: jspb.Message.getFieldWithDefault(msg, 2, 0),
    chatdatalistList: jspb.Message.toObjectList(msg.getChatdatalistList(),
    proto.model.MsgChat.toObject, includeInstance),
    chatreplylistList: jspb.Message.toObjectList(msg.getChatreplylistList(),
    proto.model.MsgChatReply.toObject, includeInstance)
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
 * @return {!proto.model.MsgChatQueryResult}
 */
proto.model.MsgChatQueryResult.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.model.MsgChatQueryResult;
  return proto.model.MsgChatQueryResult.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.model.MsgChatQueryResult} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.model.MsgChatQueryResult}
 */
proto.model.MsgChatQueryResult.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setUserid(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setToid(value);
      break;
    case 3:
      var value = new proto.model.MsgChat;
      reader.readMessage(value,proto.model.MsgChat.deserializeBinaryFromReader);
      msg.addChatdatalist(value);
      break;
    case 4:
      var value = new proto.model.MsgChatReply;
      reader.readMessage(value,proto.model.MsgChatReply.deserializeBinaryFromReader);
      msg.addChatreplylist(value);
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
proto.model.MsgChatQueryResult.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.model.MsgChatQueryResult.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.model.MsgChatQueryResult} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.model.MsgChatQueryResult.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserid();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getToid();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
  f = message.getChatdatalistList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.model.MsgChat.serializeBinaryToWriter
    );
  }
  f = message.getChatreplylistList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.model.MsgChatReply.serializeBinaryToWriter
    );
  }
};


/**
 * optional int64 userId = 1;
 * @return {number}
 */
proto.model.MsgChatQueryResult.prototype.getUserid = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.MsgChatQueryResult} returns this
 */
proto.model.MsgChatQueryResult.prototype.setUserid = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional int64 toId = 2;
 * @return {number}
 */
proto.model.MsgChatQueryResult.prototype.getToid = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.model.MsgChatQueryResult} returns this
 */
proto.model.MsgChatQueryResult.prototype.setToid = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * repeated MsgChat chatDataList = 3;
 * @return {!Array<!proto.model.MsgChat>}
 */
proto.model.MsgChatQueryResult.prototype.getChatdatalistList = function() {
  return /** @type{!Array<!proto.model.MsgChat>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.MsgChat, 3));
};


/**
 * @param {!Array<!proto.model.MsgChat>} value
 * @return {!proto.model.MsgChatQueryResult} returns this
*/
proto.model.MsgChatQueryResult.prototype.setChatdatalistList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.model.MsgChat=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.MsgChat}
 */
proto.model.MsgChatQueryResult.prototype.addChatdatalist = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.model.MsgChat, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.MsgChatQueryResult} returns this
 */
proto.model.MsgChatQueryResult.prototype.clearChatdatalistList = function() {
  return this.setChatdatalistList([]);
};


/**
 * repeated MsgChatReply chatReplyList = 4;
 * @return {!Array<!proto.model.MsgChatReply>}
 */
proto.model.MsgChatQueryResult.prototype.getChatreplylistList = function() {
  return /** @type{!Array<!proto.model.MsgChatReply>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.model.MsgChatReply, 4));
};


/**
 * @param {!Array<!proto.model.MsgChatReply>} value
 * @return {!proto.model.MsgChatQueryResult} returns this
*/
proto.model.MsgChatQueryResult.prototype.setChatreplylistList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.model.MsgChatReply=} opt_value
 * @param {number=} opt_index
 * @return {!proto.model.MsgChatReply}
 */
proto.model.MsgChatQueryResult.prototype.addChatreplylist = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.model.MsgChatReply, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.model.MsgChatQueryResult} returns this
 */
proto.model.MsgChatQueryResult.prototype.clearChatreplylistList = function() {
  return this.setChatreplylistList([]);
};


