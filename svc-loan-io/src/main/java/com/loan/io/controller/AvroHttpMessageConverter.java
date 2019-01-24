package com.loan.io.controller;

import java.io.IOException;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.List;

import org.apache.avro.Schema;
import org.apache.avro.specific.SpecificRecord;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpInputMessage;
import org.springframework.http.HttpOutputMessage;
import org.springframework.http.MediaType;
import org.springframework.http.converter.AbstractHttpMessageConverter;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.http.converter.HttpMessageNotWritableException;

public class AvroHttpMessageConverter<T> extends AbstractHttpMessageConverter<Object>
{
    protected final Log logger = LogFactory.getLog(getClass());
    
    protected String HTTP_HEADER_NAME = "json__TypeId__";

    public static final Charset DEFAULT_CHARSET = Charset.forName("UTF-8");

    /*
     * public static AvroHttpMessageConverter getMapper() { return new
     * AvroHttpMessageConverter(); }
     */

    public AvroHttpMessageConverter()
    {
        super();
        List<MediaType> supportedMediaTypes = new ArrayList<MediaType>();
         supportedMediaTypes.add( new MediaType("application", "json", DEFAULT_CHARSET));
        // supportedMediaTypes.add( new MediaType("application", "*+json",
        // DEFAULT_CHARSET));
        supportedMediaTypes.add(new MediaType("avro", "binary", DEFAULT_CHARSET));

        this.setSupportedMediaTypes(supportedMediaTypes);
    }

    protected boolean supports(Class<?> clazz)
    {
        if (SpecificRecord.class.isAssignableFrom(clazz) || String.class.isAssignableFrom(clazz))
        {
            return true;
        }
        return false;
    }

    @Override
    protected Object readInternal(Class<? extends Object> clazz, HttpInputMessage inputMessage) throws IOException, HttpMessageNotReadableException
    {

        Schema schema = AvroConverter.getSchema(clazz);

        if (null == schema)
        {
            return null;
        }
        
        return (T) AvroConverter.convertFromJson(inputMessage.getBody(), schema, clazz.getClass());
    }

    @Override
    protected void writeInternal(Object t, HttpOutputMessage outputMessage) throws IOException, HttpMessageNotWritableException
    {
        Schema schema = AvroConverter.getSchema(t.getClass());

        byte[] returnObject = AvroConverter.convertToJson(t, schema);

        HttpHeaders headers = outputMessage.getHeaders();

        headers.set(HTTP_HEADER_NAME, t.getClass().getCanonicalName().toString());

        outputMessage.getBody().write(returnObject);
    }
}