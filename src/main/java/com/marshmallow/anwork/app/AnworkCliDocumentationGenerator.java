package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.DocumentationGeneratorFactory;
import com.marshmallow.anwork.app.cli.DocumentationType;

import java.io.PrintWriter;

/**
 * This is an app that generates the documentation for the ANWORK CLI.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkCliDocumentationGenerator {

  private static final String FILENAME = "doc/CLI.md";

  /**
   * This is the main method for the {@link AnworkCliDocumentationGenerator} app.
   *
   * @param args The command line arguments for this app.
   */
  public static void main(String[] args) {
    try (PrintWriter writer = new PrintWriter(FILENAME)) {
      DocumentationGeneratorFactory.getInstance()
                                   .createGenerator(DocumentationType.GITHUB_MARKDOWN)
                                   .generate(AnworkApp.createCli(), writer);
    } catch (Exception e) {
      System.out.println("Error: " + e);
    }
  }
}
