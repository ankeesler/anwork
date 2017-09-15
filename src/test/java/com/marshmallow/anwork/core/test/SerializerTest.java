package com.marshmallow.anwork.core.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;

import com.marshmallow.anwork.core.Serializer;

/**
 * This is a generic test class for a serializable object (i.e., an object that
 * has a {@link Serializer}.
 *
 * <p>
 * Created Aug 31, 2017
 * </p>
 *
 * @author Andrew
 */
public class SerializerTest<T> {

  private Serializer<T> serializer;

  public SerializerTest(Serializer<T> serializer) {
    this.serializer = serializer;
  }

  /**
   * Assert that this string is a valid serialization of an object of type T.
   *
   * @param string The valid serialization of an object of type T.
   */
  protected T assertGood(String string) {
    T t = serializer.unmarshall(string);
    String message = String.format("Expected %s from %s", getParameterizedTypeName(), string);
    assertNotNull(message, t);

    message = String.format("Marshalled (%s) does not match original string (%s)", t, string);
    assertEquals(message, string, serializer.marshall(t));

    return t;
  }

  /**
   * Assert that this string is an invalid serialization of an object of type T.
   *
   * @param string The invalid serialization of an object of type T.
   */
  protected void assertBad(String string) {
    T t = serializer.unmarshall(string);
    String message = String.format("Expected no %s from %s", getParameterizedTypeName(), string);
    assertNull(message, t);
  }

  private String getParameterizedTypeName() {
    // TODO: how do I do this the right way?
    return "type";
  }
}
