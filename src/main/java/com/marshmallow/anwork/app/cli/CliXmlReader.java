package com.marshmallow.anwork.app.cli;

import java.io.File;
import java.net.URL;

import javax.xml.XMLConstants;
import javax.xml.parsers.ParserConfigurationException;
import javax.xml.parsers.SAXParser;
import javax.xml.parsers.SAXParserFactory;
import javax.xml.validation.Schema;
import javax.xml.validation.SchemaFactory;

import org.xml.sax.SAXException;

/**
 * An object to read an XML {@link File} and produce a {@link Cli} object.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public class CliXmlReader {

  private static final String SCHEMA_RESOURCE = "cli.xsd";

  // Use #getSchema to access me!!!
  private static Schema SCHEMA = null;

  // Accessor for static schema.
  private static Schema getSchema() throws SAXException {
    if (SCHEMA == null) {
      SchemaFactory schemaFactory = SchemaFactory.newInstance(XMLConstants.W3C_XML_SCHEMA_NS_URI);
      URL url = CliXmlReader.class.getResource(SCHEMA_RESOURCE);
      SCHEMA = schemaFactory.newSchema(url);
    }
    return SCHEMA;
  }

  private File xmlFile;

  /**
   * Create a {@link CliXmlReader} for an XML {@link File}.
   *
   * @param xmlFile The XML {@link File} to read
   */
  public CliXmlReader(File xmlFile) {
    this.xmlFile = xmlFile;
  }

  /**
   * Read a {@link Cli} object from the XML {@link File} provided to the constructor.
   *
   * @return A {@link Cli} object from the XML {@link File} provided to the constructor
   * @throws Exception if something goes wrong
   */
  public Cli read() throws Exception {
    SAXParser xmlParser = makeParser();
    CliXmlParser cliXmlParser = new CliXmlParser();
    xmlParser.parse(xmlFile, cliXmlParser);
    return cliXmlParser.getCli();
  }

  private SAXParser makeParser() throws ParserConfigurationException, SAXException {
    SAXParserFactory parserFactory = SAXParserFactory.newInstance();
    parserFactory.setSchema(getSchema());
    return parserFactory.newSAXParser();
  }
}
