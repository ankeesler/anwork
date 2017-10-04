package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliXmlReader;

import java.io.File;
import java.net.URL;

/**
 * This class creates the CLI for the ANWORK app.
 *
 * <p>
 * Created Sep 11, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkCliCreator {

  private static final String CLI_XML_RESOURCE = "anwork-cli.xml";

  /**
   * Create an instance of the CLI for the ANWORK app.
   *
   * @return An instance of the CLI for the ANWORK app
   */
  public Cli makeCli() throws Exception {
    URL xmlUrl = getClass().getResource(CLI_XML_RESOURCE);
    File xmlFile = new File(xmlUrl.toURI());
    CliXmlReader xmlReader = new CliXmlReader(xmlFile);
    return xmlReader.read();
  }
}
