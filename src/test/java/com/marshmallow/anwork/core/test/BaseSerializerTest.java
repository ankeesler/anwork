package com.marshmallow.anwork.core.test;

import static org.junit.Assert.assertNotNull;

import com.marshmallow.anwork.core.Serializable;
import com.marshmallow.anwork.core.Serializer;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;

/**
 * This is a generic test class for a serializable object.
 *
 * <p>
 * Created Aug 31, 2017
 * </p>
 *
 * <p>
 * Refactored September 16, 2017
 * </p>
 *
 * @author Andrew
 */
public class BaseSerializerTest<T extends Serializable<?>> {

  private final Serializer<T> serializer;

  /**
   * Initialize this instance with a {@link Serializer} to test.
   *
   * @param serializer The serializer that is under test
   */
  public BaseSerializerTest(Serializer<T> serializer) {
    this.serializer = serializer;
  }

  /**
   * Serialize an object and then unserialize it.
   *
   * @param t The object to serialize
   * @return The unserialized object
   * @throws IOException if any of the un/serialization fails
   */
  protected T runSerialization(T t) throws IOException {
    ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
    serializer.serialize(t, outputStream);
    InputStream inputStream = new ByteArrayInputStream(outputStream.toByteArray());
    T unserialized = serializer.unserialize(inputStream);
    assertNotNull(unserialized);
    return unserialized;
  }
}