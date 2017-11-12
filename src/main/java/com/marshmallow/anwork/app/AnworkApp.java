package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliXmlReader;

import java.io.InputStream;

/**
 * This is the main class for the anwork app.
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkApp {

  private static final String CLI_XML_RESOURCE = "anwork-cli.xml";

  /**
   * ANWORK main method.
   *
   * @param args Command line argument
   * @throws Exception for any runtime error
   */
  public static void main(String[] args) throws Exception {
    try {
      new AnworkApp().createCli().parse(args);
    } catch (Exception e) {
      throw e;
    }
  }

  /**
   * Create the CLI for the ANWORK app.
   *
   * @return The CLI for the ANWORK app.
   * @throws Exception if something goes wrong with creating the CLI
   */
  public Cli createCli() throws Exception {
    try (InputStream xmlStream = AnworkApp.class.getResourceAsStream(CLI_XML_RESOURCE)) {
      return new CliXmlReader(xmlStream).read();
    }
  }
}
