package com.marshmallow.anwork.core;

/**
 * This is a dumb factory class that just spits out blank instances of types.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Factory<T> {

  /**
   * Make a blank instance of a type T.
   *
   * @return A blank instance of a type T.
   */
  public T makeBlankInstance();
}
