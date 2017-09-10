package com.marshmallow.anwork.core;

/**
 * This is an object that can 1) turn an Object into a String and 2) given a
 * String, turn it into an Object.
 *
 * @author Andrew
 * Created Aug 31, 2017
 */
public interface Serializer<T> {

  /**
   * Turn an Object into a String.
   *
   * @param t The object to turn into a String
   * @return A string representation of the object
   */
  public String marshall(T t);

  /**
   * Turn a String into an Object.
   *
   * @param string The string to turn into an Object
   * @return An Object represented by the String
   */
  public T unmarshall(String string);
}
