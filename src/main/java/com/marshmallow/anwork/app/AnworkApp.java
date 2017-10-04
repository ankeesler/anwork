package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliXmlReader;

import java.io.File;
import java.net.URL;

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
   */
  public static void main(String[] args) {
    try {
      createCli().parse(args);
    } catch (Exception e) {
      System.out.println("Error: " + e.getMessage());
    }
  }

  /**
   * Create the CLI for the ANWORK app.
   *
   * @return The CLI for the ANWORK app.
   * @throws Exception if something goes wrong with creating the CLI
   */
  public static Cli createCli() throws Exception {
    URL xmlUrl = AnworkApp.class.getResource(CLI_XML_RESOURCE);
    File xmlFile = new File(xmlUrl.toURI());
    CliXmlReader xmlReader = new CliXmlReader(xmlFile);
    return xmlReader.read();

  }
}
