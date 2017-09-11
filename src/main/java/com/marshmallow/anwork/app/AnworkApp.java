package com.marshmallow.anwork.app;

/**
 * This is the main class for the anwork app.
 *
 * @author Andrew
 * Created Sep 9, 2017
 */
public class AnworkApp {

  /**
   * ANWORK main method.
   *
   * @param args Command line argument
   */
  public static void main(String[] args) {
    try {
      new AnworkCliCreator(new AnworkAppConfig()).makeCli().parse(args);
    } catch (Exception e) {
      System.out.println("Error: " + e.getMessage());
    }
  }
}
