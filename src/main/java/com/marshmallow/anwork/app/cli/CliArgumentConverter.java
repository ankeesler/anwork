package com.marshmallow.anwork.app.cli;

import java.util.IllegalFormatException;

/**
 * This is an object that converts from a {@link String} in the CLi system to some Java object.
 *
 * <p>
 * Note: this interface is currently private to this package. Right now, users may not
 * create their own {@CliArgumentConverter} instances. This is for simplicity. In the future, if we
 * want to offer this as a public interface, we can do so.
 * </p>
 *
 * <p>
 * Created Oct 3, 2017
 * </p>
 *
 * @author Andrew
 */
interface CliArgumentConverter<T> {

  /**
   * This is a simple method that easily allows clients to figure out to which Java object type
   * <code>this</code> {@link CliArgumentConverter} converts {@link String}'s.
   *
   * @return The Java object type to which <code>this</code> {@link CliArgumentConverter} converts
   *     {@link String}'s.
   */
  public Class<T> getConversionClass();

  /**
   * Convert a {@link String} to some Java object.
   *
   * @param string The string to convert to a Java object
   * @return A Java object that was originally a string
   * @throws IllegalFormatException if the passed {@link String} cannot be converted to the Java
   *     object type to which this class applies
   */
  public T convert(String string) throws IllegalFormatException;
}
