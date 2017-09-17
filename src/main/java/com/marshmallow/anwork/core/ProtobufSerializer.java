package com.marshmallow.anwork.core;

import com.google.protobuf.Message;
import com.google.protobuf.Parser;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;

/**
 * This is a {@link Serializer} that serializes objects to an {@link OutputStream} and from an
 * {@link InputStream} using Protocol Buffers.
 *
 * <p>
 * Created September 17, 2017
 * </p>
 *
 * @author Andrew
 */
public class ProtobufSerializer<ProtobufT extends Message,
                                SerializedT extends Serializable<ProtobufT>>
                           implements Serializer<SerializedT> {

  private Factory<SerializedT> factory;
  private Parser<ProtobufT> parser;

  /**
   * Initialize this serializer with a factory (to create new instances of type SerializedT) and a
   * {@link Parser} to parse {@link InputStream}'s of data.
   *
   * @param factory A factory to create new instances of type SerializedT
   * @param parser A parser to parse {@link InputStream}'s of data
   */
  public ProtobufSerializer(Factory<SerializedT> factory, Parser<ProtobufT> parser) {
    this.factory = factory;
    this.parser = parser;
  }

  @Override
  public void serialize(SerializedT serialized, OutputStream outputStream) throws IOException {
    ProtobufT protobuf = serialized.marshall();
    protobuf.writeDelimitedTo(outputStream);
  }

  @Override
  public SerializedT unserialize(InputStream inputStream) throws IOException {
    SerializedT serialized = factory.makeBlankInstance();
    ProtobufT protobuf = parser.parseDelimitedFrom(inputStream);
    serialized.unmarshall(protobuf);
    return serialized;
  }
}
