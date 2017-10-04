package com.marshmallow.anwork.app.cli;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;

import javax.xml.XMLConstants;
import javax.xml.parsers.ParserConfigurationException;
import javax.xml.parsers.SAXParser;
import javax.xml.parsers.SAXParserFactory;
import javax.xml.transform.stream.StreamSource;
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
  private static Schema getSchema() throws SAXException, IOException {
    if (SCHEMA == null) {
      SchemaFactory schemaFactory = SchemaFactory.newInstance(XMLConstants.W3C_XML_SCHEMA_NS_URI);
      try (InputStream schemaStream = CliXmlReader.class.getResourceAsStream(SCHEMA_RESOURCE)) {
        SCHEMA = schemaFactory.newSchema(new StreamSource(schemaStream));
      }
    }
    return SCHEMA;
  }

  private InputStream xmlStream;

  /**
   * Create a {@link CliXmlReader} with a {@link InputStream} from which to parse the CLI XML data.
   *
   * @param xmlStream The {@link InputStream} which contains the CLI XML data
   */
  public CliXmlReader(InputStream xmlStream) {
    this.xmlStream = xmlStream;
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
    xmlParser.parse(xmlStream, cliXmlParser);
    return cliXmlParser.getCli();
  }

  private SAXParser makeParser() throws ParserConfigurationException, SAXException, IOException {
    SAXParserFactory parserFactory = SAXParserFactory.newInstance();
    parserFactory.setSchema(getSchema());
    return parserFactory.newSAXParser();
  }
}
